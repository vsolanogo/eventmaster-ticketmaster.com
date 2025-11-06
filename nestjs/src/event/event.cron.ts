import { Injectable, OnModuleInit } from '@nestjs/common';
import { Cron } from '@nestjs/schedule';
import { EventService } from './event.service';

@Injectable()
export class EventCron implements OnModuleInit {
  constructor(private readonly eventService: EventService) {}

  // The first field represents minutes. In this case, it's set to 0, meaning the cron job will run at the 0th minute of the hour.
  // The second field represents hours. Here, */6 means "every 6 hours", so the job will run every 6 hours, starting from midnight (0:00), then at 6:00, 12:00, 18:00, and so on.
  // The remaining fields (day of month, month, day of week) are set to asterisks, meaning "any value". So, the cron job will run on any day of the month, in any month, and on any day of the week.

  @Cron('0 */6 * * *')
  async handleCron() {
    await this.eventService.fetchEvents();
  }

  async getTicketmastersEvent() {}

  async onModuleInit() {
    await this.eventService.fetchEvents();
  }
}
