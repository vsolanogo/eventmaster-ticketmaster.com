import { Module } from '@nestjs/common';
import { SessionModule } from './session/session.module';
import { UserModule } from './user/user.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { Session } from './session/session.entity';
import { User } from './user/user.entity';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { LoginModule } from './login/login.module';
import { RegisterModule } from './register/register.module';
import { RoleModule } from './role/role.module';
import { Role } from './role/role.entity';
import { EventModule } from './event/event.module';
import { Event } from './event/event.entity';
import { ServeStaticModule } from '@nestjs/serve-static';
import { join } from 'path';
import { ImageModule } from './image/image.module';
import { Image } from './image/image.entity';
import { ParticipantModule } from './participant/participant.module';
import { Participant } from './participant/participant.entity';
import { ScheduleModule } from '@nestjs/schedule';

@Module({
  imports: [
    ScheduleModule.forRoot(),
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, '..', 'public'),
    }),
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: `.env.${process.env.NODE_ENV}`,
    }),
    TypeOrmModule.forRootAsync({
      inject: [ConfigService],
      useFactory: (config: ConfigService) => {
        return {
          database: config.get<string>('DB_NAME'),
          // type: 'mysql',
          type: 'postgres',
          host: config.get<string>('DB_HOST'),
          port: config.get<number>('DB_PORT'),
          username: config.get<string>('DB_USERNAME'),
          password: config.get<string>('DB_PASSWORD'),
          entities: [User, Session, Role, Event, Image, Participant],
          synchronize: true,
        };
      },
    }),
    SessionModule,
    UserModule,
    LoginModule,
    RegisterModule,
    RoleModule,
    EventModule,
    ImageModule,
    ParticipantModule,
  ],
})
export class AppModule {}
