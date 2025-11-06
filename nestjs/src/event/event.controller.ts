import {
  Controller,
  UseGuards,
  Post,
  Body,
  Get,
  Query,
  Param,
  NotFoundException,
} from '@nestjs/common';
import { LoginGuard } from '../login/login.guard';
import { CurrentUser } from '../user/decorators/current-user.decorator';
import { CreateEventDto } from './dto/create-event.dto';
import { Event } from './event.entity';
import { EventService } from './event.service';
import { User } from '../user/user.entity';
import { Serialize } from '../interceptors/serialize.interceptor';
import { EventDto } from './dto/event.dto';
import { EventListDto } from './dto/EventListDto.dto';

@Controller('events')
export class EventController {
  constructor(private readonly eventService: EventService) {}

  @Serialize(EventDto)
  @Post()
  @UseGuards(LoginGuard)
  async create(
    @Body() body: CreateEventDto,
    @CurrentUser() user: User,
  ): Promise<Event> {
    const event = await this.eventService.create(body, user);
    return event;
  }

  @Serialize(EventListDto)
  @Get()
  async getAllEvents(
    @Query('page') page: string = '1',
    @Query('limit') limit: string = '10',
    @Query('sortBy') sortBy: string = 'eventDate',
    @Query('sortOrder') sortOrder: 'ASC' | 'DESC' = 'ASC',
  ): Promise<{ events: Event[]; totalCount: number }> {
    const parsedPage = parseInt(page, 10);
    const parsedLimit = parseInt(limit, 10);
    const event = await this.eventService.getPaginatedEvents(
      parsedPage,
      parsedLimit,
      sortBy,
      sortOrder,
    );
    return event;
  }

  @Serialize(EventDto)
  @Get(':id')
  async getEventById(@Param('id') id: string): Promise<Event> {
    const event = await this.eventService.getEventById(id);
    if (!event) {
      throw new NotFoundException('Event not found');
    }
    return event;
  }
}
