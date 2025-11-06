import { Module } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { RegisterController } from './register.controller';
import { UserService } from '../user/user.service';
import { User } from '../user/user.entity';
import { Role } from '../role/role.entity';

@Module({
  imports: [TypeOrmModule.forFeature([User, Role])],
  controllers: [RegisterController],
  providers: [UserService],
})
export class RegisterModule {}
