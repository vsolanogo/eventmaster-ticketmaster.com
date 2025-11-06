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
import { IsString, IsEmail } from 'class-validator';
import { Role } from '../role/role.entity';

@Entity()
export class User {
  @PrimaryGeneratedColumn('uuid')
  id: string;

  @OneToMany(() => Session, (i) => i.user)
  session: Session[];

  @ManyToMany(() => Role)
  @JoinTable()
  role: Role[];

  @Column({ nullable: false })
  password: string;

  @Column({ nullable: false, length: 255, unique: true })
  @IsString()
  @IsEmail()
  email: string;

  @CreateDateColumn()
  createdAt: Date;

  @UpdateDateColumn()
  updatedAt: Date;
}
