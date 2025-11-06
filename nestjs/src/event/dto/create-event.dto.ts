import {
  IsNotEmpty,
  IsString,
  IsNumber,
  Min,
  IsOptional,
  MinLength,
  MaxLength,
  IsArray,
  Max,
  IsDateString,
} from 'class-validator';

export class CreateEventDto {
  @IsNotEmpty()
  @IsString()
  @MinLength(2, { message: 'Title is too short' })
  title: string;

  @IsOptional()
  @IsArray()
  images: string[]; // Define the new field as an array of strings

  @IsString()
  @MaxLength(5000, { message: 'Description is too long' })
  description: string;

  @IsDateString()
  eventDate: Date;

  @IsOptional()
  @IsNumber()
  @Min(-90)
  @Max(90)
  latitude: number;

  @IsOptional()
  @IsNumber()
  @Min(-180)
  @Max(180)
  longitude: number;
}
