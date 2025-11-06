import { Event } from '../event/event.entity';
import { Entity, PrimaryGeneratedColumn, Column, ManyToOne } from 'typeorm';
import { IsString, IsNotEmpty } from 'class-validator';

@Entity()
export class Image {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  @IsString()
  @IsNotEmpty()
  link: string;

  @ManyToOne(() => Event, (i) => i.images)
  event: Event;
}
