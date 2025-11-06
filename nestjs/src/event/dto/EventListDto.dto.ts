import { Expose, Type } from 'class-transformer';
import { EventDto } from './event.dto';

export class EventListDto {
  @Expose()
  @Type(() => EventDto)
  events: EventDto[];

  @Expose()
  totalCount: number;
}
