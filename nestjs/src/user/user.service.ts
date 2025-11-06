import { Injectable, BadRequestException, Logger } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { User } from './user.entity';
import { Repository } from 'typeorm';
import { RegisterDto } from './dto/user-register.dto';
import { hashPassword } from 'metautil';
import { validate } from 'class-validator';
import { Role } from '../role/role.entity';
import { RolesEnum } from '../models/models';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class UserService {
  private readonly adminEmail: string;
  private readonly adminPassword: string;

  constructor(
    @InjectRepository(User)
    private readonly userRepository: Repository<User>,
    @InjectRepository(Role)
    private readonly roleRepository: Repository<Role>,
    private readonly configService: ConfigService,
  ) {
    this.adminEmail = this.configService.get<string>('ROOT_ADMIN_EMAIL');
    this.adminPassword = this.configService.get<string>('ROOT_ADMIN_PASSWORD');
  }

  async create(registerDto: RegisterDto): Promise<User> {
    const existingUser = await this.userRepository.findOne({
      where: [{ email: registerDto.email }],
    });

    if (existingUser) {
      throw new BadRequestException('User already exists.');
    }

    // by default each user will be registered with user role
    const role = await this.roleRepository.findOne({
      where: { role: RolesEnum.User },
    });

    const newUser = new User();
    newUser.email = registerDto.email;
    const hash = await hashPassword(registerDto.password);
    newUser.password = hash;
    newUser.role = [role];

    const errors = await validate(newUser);
    if (errors.length > 0) {
      throw new BadRequestException(errors);
    } else {
      return this.userRepository.save(newUser);
    }
  }

  async seed(): Promise<void> {
    const existingAdmin = await this.userRepository.findOne({
      where: { email: this.adminEmail },
    });

    if (existingAdmin) {
      return;
    }

    const role = await this.roleRepository.findOne({
      where: { role: RolesEnum.Admin },
    });

    const adminUser = new User();
    adminUser.email = this.adminEmail;
    adminUser.role = [role];
    const hash = await hashPassword(this.adminPassword);
    adminUser.password = hash;

    const errors = await validate(adminUser);
    if (errors.length > 0) {
      throw new BadRequestException(errors);
    } else {
      this.userRepository.save(adminUser);
    }
  }
}
