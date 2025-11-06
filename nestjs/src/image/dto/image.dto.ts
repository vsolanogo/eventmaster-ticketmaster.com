import { Expose, Type } from 'class-transformer';
import { EventDto } from 'src/event/dto/event.dto';

export class ImageDto {
  @Expose()
  id: string;

  @Expose()
  link: string;

  @Expose()
  @Type(() => EventDto)
  event: EventDto;
}
