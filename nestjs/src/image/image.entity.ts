import { Event } from '../event/event.entity';
import { Entity, PrimaryGeneratedColumn, Column, ManyToOne } from 'typeorm';

@Entity()
export class Image {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  link: string;

  @ManyToOne(() => Event, (i) => i.images)
  event: Event;
}
