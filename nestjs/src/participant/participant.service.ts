import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Event } from '../event/event.entity';
import { validate } from 'class-validator';
import { Participant } from './participant.entity';
import { CreateParticipantDto } from './dto/create-participant.dto';
import { randEmail, randPassword, randFullName } from '@ngneat/falso';
import {
  getRandomDate,
  getRandomSourceOfEventDiscovery,
} from 'src/urils/random.util';

@Injectable()
export class ParticipantService {
  constructor(
    @InjectRepository(Participant)
    private readonly participantRepository: Repository<Participant>,
    @InjectRepository(Event)
    private readonly eventRepository: Repository<Event>,
  ) {}

  async create(
    createParticipantDto: CreateParticipantDto,
    eventId: string,
  ): Promise<Participant> {
    const event = await this.eventRepository.findOne({
      where: { id: eventId },
    });

    if (!event) {
      throw new NotFoundException('Event not found');
    }

    const newParticipant = new Participant();
    newParticipant.email = createParticipantDto.email;
    newParticipant.dateOfBirth = createParticipantDto.dateOfBirth;
    newParticipant.fullName = createParticipantDto.fullName;
    newParticipant.sourceOfEventDiscovery =
      createParticipantDto.sourceOfEventDiscovery;
    newParticipant.dateOfBirth = new Date(createParticipantDto.dateOfBirth);
    newParticipant.event = event;

    const errors = await validate(newParticipant);

    if (errors.length > 0) {
      throw new BadRequestException(errors);
    } else {
      return this.participantRepository.save(newParticipant);
    }
  }

  async generateFakeParticipants(event: Event, count: number): Promise<void> {
    const fakeParticipants: Participant[] = [];
    for (let i = 0; i < count; i++) {
      const fakeParticipant = new Participant();
      fakeParticipant.fullName = randFullName();
      fakeParticipant.email = randEmail();
      fakeParticipant.dateOfBirth = new Date(1990, 0, 1); // Fixed date for simplicity
      fakeParticipant.sourceOfEventDiscovery =
        getRandomSourceOfEventDiscovery();
      fakeParticipant.event = event;
      fakeParticipant.createdAt = getRandomDate(
        new Date(2023, 0, 1),
        new Date(),
      ); // Set random createdAt date
      fakeParticipants.push(fakeParticipant);
    }
    await this.participantRepository.save(fakeParticipants);
  }

  async getAllByEventId(eventId: string): Promise<Participant[]> {
    const event = await this.eventRepository.findOne({
      where: { id: eventId },
    });

    if (!event) {
      throw new NotFoundException('Event not found');
    }

    return this.participantRepository.find({
      where: { event: { id: eventId } },
    });
  }

  async getRegistrationsPerDay(
    eventId: string,
  ): Promise<{ date: string; count: number }[]> {
    const event = await this.eventRepository.findOne({
      where: { id: eventId },
    });

    if (!event) {
      throw new NotFoundException('Event not found');
    }

    const registrations = await this.participantRepository
      .createQueryBuilder('participant')
      .select('DATE(participant.createdAt)', 'date')
      .addSelect('COUNT(participant.id)', 'count')
      .where('participant.eventId = :eventId', { eventId })
      .groupBy('DATE(participant.createdAt)')
      .orderBy('date')
      .getRawMany();

    return registrations.map(({ date, count }) => ({
      date,
      count: parseInt(count, 10),
    }));
  }
}
