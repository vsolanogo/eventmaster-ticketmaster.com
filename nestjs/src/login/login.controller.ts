import {
  Controller,
  Post,
  Body,
  HttpStatus,
  Ip,
  Response,
} from '@nestjs/common';
import { LoginService } from './login.service';
import { LoginDto } from './dto/login.dto';
import { SESSION_ID } from '../constants';

@Controller()
export class LoginController {
  constructor(private readonly loginService: LoginService) {}

  @Post('login')
  async auth(@Body() body: LoginDto, @Ip() ip, @Response() res): Promise<void> {
    const newSession = await this.loginService.auth(body, ip);

    res.cookie(SESSION_ID, newSession.token, {
      expires: newSession.expires,
      sameSite: 'strict',
      httpOnly: true,
    });

    return res.status(HttpStatus.OK).send();
  }

  @Post('logout')
  async logout(@Response() res): Promise<void> {
    res.clearCookie(SESSION_ID);
    // TODO: remove token from db here

    return res.status(HttpStatus.OK).send();
  }
}
