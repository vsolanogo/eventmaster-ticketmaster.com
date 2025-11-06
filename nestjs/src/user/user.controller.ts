import { Controller, Req, Get, UseGuards } from '@nestjs/common';
import { SESSION_ID } from '../constants';
import { SessionService } from '../session/session.service';
import { LoginGuard } from '../login/login.guard';
import { Serialize } from '../interceptors/serialize.interceptor';
import { UserDto } from './dto/user.dto';
import { User } from './user.entity';
import { CurrentUser } from './decorators/current-user.decorator';

@Controller('user')
export class UserController {
  constructor(private readonly sessionService: SessionService) {}

  @Serialize(UserDto)
  @Get()
  @UseGuards(LoginGuard)
  async get(@CurrentUser() user: User): Promise<User> {
    return user;
  }
}
