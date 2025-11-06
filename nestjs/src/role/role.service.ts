import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { Role } from './role.entity';
import { RolesEnum } from '../models/models';

@Injectable()
export class RoleService {
  constructor(
    @InjectRepository(Role)
    private readonly roleRepository: Repository<Role>,
  ) {}

  async seed(): Promise<void> {
    const existingRoles = await this.roleRepository.find({
      where: [{ role: RolesEnum.Admin }, { role: RolesEnum.User }],
    });

    const rolesToSave: Role[] = [];

    if (!existingRoles.find((role) => role.role === RolesEnum.Admin)) {
      const admin = new Role();
      admin.role = RolesEnum.Admin;
      admin.description = null;
      rolesToSave.push(admin);
    }

    if (!existingRoles.find((role) => role.role === RolesEnum.User)) {
      const user = new Role();
      user.role = RolesEnum.User;
      user.description = null;
      rolesToSave.push(user);
    }

    await this.roleRepository.save(rolesToSave);
  }
}
