import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  CreateDateColumn,
  UpdateDateColumn,
} from 'typeorm';
import { User } from '../user/user.entity';
import {
  IsString,
  IsNotEmpty,
  IsDate,
  ValidateNested,
  IsOptional,
} from 'class-validator';
import { Type } from 'class-transformer';

@Entity()
export class Session {
  @PrimaryGeneratedColumn('uuid')
  @IsOptional()
  id: string;

  @ManyToOne(() => User, (i) => i.session, { nullable: false })
  @ValidateNested()
  @Type(() => User)
  user: User;

  @Column({ nullable: false, length: 128, unique: true })
  @IsString()
  @IsNotEmpty()
  token: string;

  @Column({ nullable: false, length: 45 })
  @IsString()
  @IsNotEmpty()
  ip: string;

  @CreateDateColumn()
  @IsOptional()
  @IsDate()
  createdAt: Date;

  @UpdateDateColumn()
  @IsOptional()
  @IsDate()
  updatedAt: Date;

  @Column()
  @IsDate()
  expires: Date;
}
