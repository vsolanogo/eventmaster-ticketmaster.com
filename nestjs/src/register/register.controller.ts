import { Controller, Body, Post, HttpCode } from '@nestjs/common';
import { RegisterDto } from '../user/dto/user-register.dto';
import { UserService } from '../user/user.service';
import { Serialize } from '../interceptors/serialize.interceptor';
import { UserDto } from '../user/dto/user.dto';
import { User } from '../user/user.entity';

@Controller('register')
export class RegisterController {
  constructor(private readonly userService: UserService) {}

  @Serialize(UserDto)
  @Post()
  @HttpCode(201)
  async register(@Body() body: RegisterDto): Promise<User> {
    const user = await this.userService.create(body);

    return user;
  }
}
