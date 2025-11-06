import { Expose, Type } from 'class-transformer';
import { RoleDto } from '../../role/dto/role.dto';
import { SessionDto } from '../../session/dto/session.dto';

export class UserDto {
  @Expose()
  id: string;

  @Expose()
  email: string;

  @Expose()
  @Type(() => RoleDto)
  role: RoleDto[];

  //do not expose on serialization
  password: string;

  @Expose()
  @Type(() => SessionDto)
  session: SessionDto[];

  @Expose()
  createdAt: Date;

  @Expose()
  updatedAt: Date;
}
