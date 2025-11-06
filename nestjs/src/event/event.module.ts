import { Module } from '@nestjs/common';
import { EventController } from './event.controller';
import { EventService } from './event.service';
import { Event } from './event.entity';
import { User } from '../user/user.entity';
import { Image } from '../image/image.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { LoginGuard } from '../login/login.guard';
import { SessionService } from '../session/session.service';
import { Session } from '../session/session.entity';
import { Participant } from '../participant/participant.entity';
import { EventCron } from './event.cron';
import { ImageService } from '../image/image.service';
import { ParticipantService } from 'src/participant/participant.service';

@Module({
  imports: [
    TypeOrmModule.forFeature([Event, User, Session, Image, Participant]),
  ],
  controllers: [EventController],
  providers: [
    EventService,
    LoginGuard,
    SessionService,
    EventCron,
    ImageService,
    ParticipantService,
  ],
})
export class EventModule {}
