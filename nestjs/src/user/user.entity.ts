import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  OneToMany,
  ManyToMany,
  JoinTable,
} from 'typeorm';
import { Session } from '../session/session.entity';
import { IsString, IsEmail, IsOptional, IsDate } from 'class-validator';
import { Role } from '../role/role.entity';

@Entity()
export class User {
  @PrimaryGeneratedColumn('uuid')
  @IsOptional()
  id: string;

  @OneToMany(() => Session, (i) => i.user)
  @IsOptional()
  session: Session[];

  @ManyToMany(() => Role)
  @JoinTable()
  @IsOptional()
  role: Role[];

  @Column({ nullable: false })
  @IsString()
  password: string;

  @Column({ nullable: false, length: 255, unique: true })
  @IsString()
  @IsEmail()
  email: string;

  @CreateDateColumn()
  @IsOptional()
  @IsDate()
  createdAt: Date;

  @UpdateDateColumn()
  @IsOptional()
  @IsDate()
  updatedAt: Date;
}
