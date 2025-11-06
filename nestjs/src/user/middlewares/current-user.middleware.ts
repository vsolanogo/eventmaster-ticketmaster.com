import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response, NextFunction } from 'express';
import { User } from '../user.entity';
import { SESSION_ID } from '../../constants';
import { SessionService } from '../../session/session.service';

declare global {
  // eslint-disable-next-line @typescript-eslint/no-namespace
  namespace Express {
    interface Request {
      currentUser?: User;
    }
  }
}

@Injectable()
export class CurrentUserMiddleware implements NestMiddleware {
  constructor(private sessionService: SessionService) {}

  async use(req: Request, res: Response, next: NextFunction) {
    const token = req.cookies?.[SESSION_ID] || null;

    if (token) {
      const user = await this.sessionService.getUserByToken(token, res);

      req.currentUser = user;
    }

    next();
  }
}
