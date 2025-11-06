import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { LoginController } from './login.controller';
import { LoginService } from './login.service';
import { Session } from '../session/session.entity';
import { User } from '../user/user.entity';
import { LoginGuard } from './login.guard';
import { SessionService } from '../session/session.service';

@Module({
  imports: [TypeOrmModule.forFeature([Session, User])],
  controllers: [LoginController],
  providers: [LoginService, LoginGuard, SessionService],
})
export class LoginModule {}
