import { Expose } from 'class-transformer';

export class RoleDto {
  @Expose()
  role: string;

  @Expose()
  description: string;
}
