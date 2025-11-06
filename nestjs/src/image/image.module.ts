import { Module } from '@nestjs/common';
import { ImageController } from './image.controller';
import { ImageService } from './image.service';
import { Image } from './image.entity';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Event } from '../event/event.entity';
import { Role } from '../role/role.entity';
import { RoleModule } from '../role/role.module';
import { Session } from '../session/session.entity';

@Module({
  imports: [TypeOrmModule.forFeature([Event, Session, Role, Image]), RoleModule],
  controllers: [ImageController],
  providers: [ImageService],
})
export class ImageModule {}
