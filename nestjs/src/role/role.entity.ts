import { Entity, Column, PrimaryColumn } from 'typeorm';
import { IsString } from 'class-validator';

@Entity()
export class Role {
  @PrimaryColumn()
  @IsString()
  role: string;

  @Column({ nullable: true })
  @IsString()
  description: string;
}
