import { Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import { User } from '../user/user.entity';
import { SessionService } from '../session/session.service';
import { LoginDto } from './dto/login.dto';
import { Session } from '../session/session.entity';
import { validatePassword } from 'metautil';
import { Injectable, UnauthorizedException } from '@nestjs/common';

@Injectable()
export class LoginService {
  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
    private readonly sessionService: SessionService,
  ) {}

  async auth(loginDto: LoginDto, ip): Promise<Session> {
    const user = await this.userRepository.findOne({
      where: { email: loginDto.email },
    });

    if (!user || !(await validatePassword(loginDto.password, user.password))) {
      throw new UnauthorizedException('Invalid credentials');
    }

    return this.sessionService.create(user.id, ip);
  }
}
