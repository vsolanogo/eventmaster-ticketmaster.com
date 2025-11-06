import { MiddlewareConsumer, Module, OnModuleInit } from '@nestjs/common';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserService } from './user.service';
import { User } from './user.entity';
import { Session } from '../session/session.entity';
import { UserController } from './user.controller';
import { SessionService } from '../session/session.service';
import { Role } from '../role/role.entity';
import { RoleModule } from '../role/role.module';
import { CurrentUserMiddleware } from './middlewares/current-user.middleware';

@Module({
  imports: [TypeOrmModule.forFeature([User, Session, Role]), RoleModule],
  providers: [UserService, SessionService],
  controllers: [UserController],
})
export class UserModule implements OnModuleInit {
  constructor(private readonly userService: UserService) {}

  configure(consumer: MiddlewareConsumer) {
    consumer.apply(CurrentUserMiddleware).forRoutes('*');
  }

  async onModuleInit() {
    await this.userService.seed();
  }
}
