import { Event } from '../event/event.entity';
import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  CreateDateColumn,
} from 'typeorm';
import { IsString, IsNotEmpty, IsEmail, IsDate, IsOptional, ValidateNested } from 'class-validator';
import { Type } from 'class-transformer';

@Entity()
export class Participant {
  @PrimaryGeneratedColumn('uuid')
  @IsOptional()
  id: string;

  @Column()
  @IsString()
  @IsNotEmpty()
  fullName: string;

  @Column()
  @IsEmail()
  @IsNotEmpty()
  email: string;

  @Column()
  @IsDate()
  @Type(() => Date)
  dateOfBirth: Date;

  @Column()
  @IsString()
  @IsNotEmpty()
  sourceOfEventDiscovery: string;

  @ManyToOne(() => Event, (i) => i.participants)
  @ValidateNested()
  @IsOptional()
  event: Event;

  @CreateDateColumn()
  @IsOptional()
  @IsDate()
  createdAt: Date;
}
