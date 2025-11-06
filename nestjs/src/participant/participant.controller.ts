import { Controller, Post, Body, Get, Param } from '@nestjs/common';
import { ParticipantService } from './participant.service';
import { CreateParticipantDto } from './dto/create-participant.dto';
import { Participant } from './participant.entity';

@Controller('participant')
export class ParticipantController {
  constructor(private readonly participantService: ParticipantService) {}

  @Get('event/:eventId/registrations-per-day')
  async getRegistrationsPerDay(
    @Param('eventId') eventId: string,
  ): Promise<{ date: string; count: number }[]> {
    return await this.participantService.getRegistrationsPerDay(eventId);
  }

  @Post('event/:eventId')
  async create(
    @Body() body: CreateParticipantDto,
    @Param('eventId') eventId: string,
  ): Promise<Participant> {
    return await this.participantService.create(body, eventId);
  }

  @Get('event/:eventId')
  async getAllByEventId(
    @Param('eventId') eventId: string,
  ): Promise<Participant[]> {
    return await this.participantService.getAllByEventId(eventId);
  }
}
