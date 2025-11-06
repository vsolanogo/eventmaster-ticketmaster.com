import {
  IsIn,
  IsString,
  IsNotEmpty,
  MinLength,
  MaxLength,
  IsEmail,
  IsDateString,
} from 'class-validator';

export class CreateParticipantDto {
  @IsNotEmpty()
  @IsString()
  @MinLength(1, { message: 'Full name is too short' })
  @MaxLength(255, { message: 'Full name is too long' })
  fullName: string;

  @IsString()
  @IsNotEmpty()
  @IsEmail()
  email: string;

  @IsDateString()
  dateOfBirth: Date;

  @IsString()
  @IsIn(['Social media', 'Friends', 'Found myself'])
  sourceOfEventDiscovery: 'Social media' | 'Friends' | 'Found myself';
}
