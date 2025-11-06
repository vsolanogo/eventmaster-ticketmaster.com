import { Expose, Type } from 'class-transformer';
import { UserDto } from '../../user/dto/user.dto';

export class SessionDto {
  @Expose()
  id: number;

  user: UserDto;

  token: string;

  @Expose()
  ip: string;

  @Expose()
  createdAt: Date;

  expires: Date;
}
