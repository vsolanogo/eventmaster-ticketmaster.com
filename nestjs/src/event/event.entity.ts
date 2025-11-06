import {
  Entity,
  Column,
  CreateDateColumn,
  UpdateDateColumn,
  ManyToOne,
  OneToMany,
  PrimaryColumn,
} from 'typeorm';
import { IsString, MaxLength, MinLength } from 'class-validator';
import { Image } from '../image/image.entity';
import { User } from '../user/user.entity';
import { Participant } from '../participant/participant.entity';

@Entity()
export class Event {
  @PrimaryColumn()
  id: string;

  @ManyToOne(() => User, (i) => i.id)
  user: User;

  @MinLength(2, { message: 'Title is too short' })
  @MaxLength(255, { message: 'Title is too long' })
  @Column({ nullable: false, length: 255 })
  title: string;

  @Column({ type: 'text', nullable: false })
  @IsString()
  @MaxLength(5000, { message: 'Description is too long' })
  description: string;

  @Column({ type: 'text' })
  organizer: string;

  @OneToMany(() => Image, (i) => i.event)
  images: Image[];

  @OneToMany(() => Participant, (i) => i.event)
  participants: Participant[];

  @Column({ type: 'decimal', precision: 10, scale: 8 })
  latitude: number;

  @Column({ type: 'decimal', precision: 11, scale: 8 })
  longitude: number;

  @Column()
  eventDate: Date;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
