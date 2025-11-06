import {
  BadRequestException,
  Injectable,
  NotFoundException,
  Logger,
} from '@nestjs/common';
import { v4 as uuidv4 } from 'uuid';
import { InjectRepository } from '@nestjs/typeorm';
import { Event } from './event.entity';
import { Image } from '../image/image.entity';
import { Repository, In } from 'typeorm';
import { CreateEventDto } from './dto/create-event.dto';
import { User } from '../user/user.entity';
import { validate } from 'class-validator';
import { ConfigService } from '@nestjs/config';
import { Client } from 'undici';
import { ImageService } from '../image/image.service';
import { FetchedEvent } from './dto/fetched-event.dto';
import { generateEventDescription } from './fetchedEventDescriptionGeneration';
import { ParticipantService } from 'src/participant/participant.service';
import { getRandomNumber } from 'src/urils/random.util';

@Injectable()
export class EventService {
  private readonly logger = new Logger(EventService.name);
  private readonly tickermasterKey: string;

  constructor(
    private readonly imageService: ImageService,
    private configService: ConfigService,
    @InjectRepository(Event)
    private readonly eventRepository: Repository<Event>,
    @InjectRepository(Image)
    private readonly imageRepository: Repository<Image>,
    private readonly participantService: ParticipantService,
  ) {
    this.tickermasterKey = this.configService.get<string>('TICKETMASTER_KEY');
  }

  async create(
    createEventDto: CreateEventDto | FetchedEvent,
    currentUser: User | null,
  ): Promise<Event | null> {
    if ('id' in createEventDto) {
      const existingEvent = await this.eventRepository.findOne({
        where: { id: createEventDto.id },
      });
      if (existingEvent) {
        // If the event exists, silently skip and return null
        return null;
      }
    }

    const newEvent = new Event();

    const images = await this.imageRepository.findBy({
      id: In(createEventDto.images),
    });
    if (createEventDto.images.length !== images.length) {
      throw new NotFoundException('Not all requested images were found');
    }

    newEvent.id = 'id' in createEventDto ? createEventDto.id : uuidv4();
    newEvent.images = images;
    newEvent.user = currentUser;
    newEvent.title = createEventDto.title;
    newEvent.description = createEventDto.description;
    newEvent.organizer =
      'organizer' in createEventDto
        ? createEventDto.organizer
        : currentUser.email;
    newEvent.latitude = createEventDto.latitude;
    newEvent.longitude = createEventDto.longitude;
    newEvent.eventDate = new Date(createEventDto.eventDate);

    const errors = await validate(newEvent);

    if (errors.length > 0) {
      throw new BadRequestException(errors);
    } else {
      const savedEvent = await this.eventRepository.save(newEvent);

      // Add fake participants if the event is FetchedEvent
      if ('id' in createEventDto) {
        await this.participantService.generateFakeParticipants(
          savedEvent,
          getRandomNumber(20, 100),
        ); // Add fake participants
      }

      return savedEvent;
    }
  }

  async findAll(): Promise<Event[]> {
    return this.eventRepository.find({
      relations: ['images', 'user'],
      order: { eventDate: 'ASC' },
    });
  }

  async getPaginatedEvents(
    page: number,
    limit: number,
    sortBy: string,
    sortOrder: 'ASC' | 'DESC',
  ): Promise<{ events: Event[]; totalCount: number }> {
    const orderOptions: Record<string, 'ASC' | 'DESC'> = {
      ASC: 'ASC',
      DESC: 'DESC',
    };

    const order = orderOptions[sortOrder];

    const sortOptions: Record<string, string> = {
      title: 'title',
      eventDate: 'eventDate',
      organizer: 'organizer',
    };

    const eventsQuery = this.eventRepository
      .createQueryBuilder('event')
      .leftJoinAndSelect('event.user', 'user')
      .leftJoinAndSelect('event.images', 'images')
      .orderBy(`event.${sortOptions[sortBy]}`, order)
      .skip((page - 1) * limit)
      .take(limit);

    const [events, totalCount] = await eventsQuery.getManyAndCount();

    return { events, totalCount };
  }

  async getEventById(id: string): Promise<Event> {
    return this.eventRepository.findOne({
      where: { id },
      relations: ['images', 'user', 'participants'],
    });
  }

  async fetchEvents() {
    const client = new Client('https://app.ticketmaster.com');

    try {
      const { body } = await client.request({
        path: `/discovery/v2/events.json?countryCode=US&size=100&apikey=${this.tickermasterKey}`,
        method: 'GET',
      });
      const data = await body.json();

      if ((data as any)._embedded && (data as any)._embedded.events) {
        const events = (data as any)._embedded.events;
        await this.saveEvents(events);
      } else {
        this.logger.warn('No events found');
      }
    } catch (error) {
      console.error('Error fetching events from Ticketmaster:', error);
    } finally {
      await client.close();
    }
  }

  private async saveEvents(events: any[]) {
    for (const event of events) {
      const eventDateTime = event.dates.start.dateTime;
      let parsedDateTime = new Date(eventDateTime);

      if (isNaN(parsedDateTime.getTime())) {
        this.logger.warn(
          `Invalid event dateTime: ${eventDateTime} for event ID: ${event.id}. Using current date.`,
        );
        parsedDateTime = new Date(); // Fallback to current date
      }

      const attrImages = event._embedded.attractions.map((i) =>
        i.images.reduce((prev, current) =>
          prev.height > current.height ? prev : current,
        ),
      );

      const images = attrImages.map((i) => i.url);

      const createdImages =
        await this.imageService.createImagesWithLinks(images);

      const eventDto: FetchedEvent = {
        id: event.id,
        title: event.name,
        images: createdImages.map((i) => i.id),
        description: generateEventDescription(event),
        organizer: event?._embedded?.venues?.[0]?.name || event.name,
        eventDate: parsedDateTime,
        latitude: parseFloat(event._embedded.venues[0].location.latitude),
        longitude: parseFloat(event._embedded.venues[0].location.longitude),
      };

      await this.create(eventDto, null);
    }
  }
}
