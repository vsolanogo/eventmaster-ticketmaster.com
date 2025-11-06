import { Module } from '@nestjs/common';
import { ParticipantController } from './participant.controller';
import { ParticipantService } from './participant.service';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Participant } from './participant.entity';
import { Event } from '../event/event.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Event, Participant])],
  controllers: [ParticipantController],
  providers: [ParticipantService],
})
export class ParticipantModule {}
