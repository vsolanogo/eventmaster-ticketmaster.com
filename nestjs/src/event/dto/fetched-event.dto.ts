import { CreateEventDto } from './create-event.dto';
import { IsString } from 'class-validator';

export class FetchedEvent extends CreateEventDto {
  @IsString()
  id: string;

  @IsString()
  organizer: string;
}
