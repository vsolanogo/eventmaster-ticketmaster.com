import {
  Bind,
  Controller,
  Post,
  UploadedFile,
  UseInterceptors,
  UseGuards,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import { LoginGuard } from '../login/login.guard';
import { ImageService } from './image.service';
import { CurrentUser } from '../user/decorators/current-user.decorator';
import { User } from '../user/user.entity';
import { Image } from './image.entity';

@Controller('image')
@UseGuards(LoginGuard)
export class ImageController {
  constructor(private readonly imageService: ImageService) {}

  @Post()
  @UseInterceptors(FileInterceptor('file'))
  @Bind(UploadedFile())
  async uploadFile(
    @UploadedFile() file,
    @CurrentUser() user: User,
  ): Promise<Image> {
    return this.imageService.create(file);
  }
}
