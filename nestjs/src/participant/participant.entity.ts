import { Event } from '../event/event.entity';
import {
  Entity,
  PrimaryGeneratedColumn,
  Column,
  ManyToOne,
  CreateDateColumn,
} from 'typeorm';

@Entity()
export class Participant {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  fullName: string;

  @Column()
  email: string;

  @Column()
  dateOfBirth: Date;

  @Column()
  sourceOfEventDiscovery: string;

  @ManyToOne(() => Event, (i) => i.participants)
  event: Event;

  @CreateDateColumn()
  createdAt: Date;
}
