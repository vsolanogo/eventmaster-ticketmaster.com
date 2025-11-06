import { Repository } from 'typeorm';
import { InjectRepository } from '@nestjs/typeorm';
import { Session } from './session.entity';
import { User } from '../user/user.entity';
import {
  BadRequestException,
  Injectable,
  NotFoundException,
} from '@nestjs/common';
import { validate } from 'class-validator';
import { generateToken } from 'metautil';
import {
  configCharacters,
  configLength,
  configSecret,
} from '../config/sessions';
import { SESSION_ID } from 'src/constants';

@Injectable()
export class SessionService {
  constructor(
    @InjectRepository(Session)
    private readonly sessionRepository: Repository<Session>,
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
  ) {}

  async create(userId: string, ip: string): Promise<Session> {
    const user = await this.userRepository.findOne({ where: { id: userId } });

    if (!user) {
      throw new NotFoundException('User not found');
    }

    const cookieExpires = new Date(
      new Date().getTime() + 7 * 24 * 60 * 60 * 1000,
    );

    const newToken = generateToken(
      configSecret,
      configCharacters,
      configLength,
    );
    const session = new Session();
    session.user = user;
    session.token = newToken;
    session.ip = ip;
    session.expires = cookieExpires;

    const errors = await validate(session);
    if (errors.length > 0) {
      throw new BadRequestException(errors);
    } else {
      return this.sessionRepository.save(session);
    }
  }

  async getUserByToken(token: string, res): Promise<User> {
    const session = await this.sessionRepository.findOne({
      relations: ['user'],

      where: { token: token },
    });

    if (!session) {
      res.clearCookie(SESSION_ID);
      throw new NotFoundException('Session not found');
    }

    const user = await this.userRepository.findOne({
      relations: ['session', 'role'],
      where: { id: session.user.id },
    });

    if (!user) {
      throw new NotFoundException('User not found');
    }

    return user;
  }
}
