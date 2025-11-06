import { BadRequestException, Injectable } from '@nestjs/common';
import { Image } from './image.entity';
import { DataSource } from 'typeorm';
import { validate } from 'class-validator';
import * as fs from 'fs';

const folder = `./public`;

@Injectable()
export class ImageService {
  constructor(private readonly dataSource: DataSource) {}

  async createImagesWithLinks(links: string[]): Promise<Image[]> {
    const createdImages: Image[] = [];
    for (const link of links) {
      // Skip invalid or empty links
      if (!link || typeof link !== 'string' || link.trim() === '') {
        continue;
      }
      try {
        const image = await this.createWithLink(link);
        createdImages.push(image);
      } catch (error) {
        console.error(
          `Failed to create image with link: ${link}`,
          error.message,
        );
        // Continue with other images instead of failing completely
      }
    }
    return createdImages;
  }

  async createWithLink(link: string): Promise<Image> {
    const newImage = new Image();
    newImage.link = link;

    const errors = await validate(newImage);

    if (errors.length > 0) {
      throw new BadRequestException(errors);
    }

    const res = await this.dataSource.manager.save(newImage);
    return res;
  }

  async create(file): Promise<Image> {
    const queryRunner = this.dataSource.createQueryRunner();
    await queryRunner.connect();
    await queryRunner.startTransaction();

    try {
      const parts = file.originalname.split('.');
      const extension = parts[parts.length - 1];

      const newImage = new Image();
      // Assign a temporary value to the link field
      newImage.link = 'temp';

      const errors = await validate(newImage);

      if (errors.length > 0) {
        throw new BadRequestException(errors);
      }

      const res = await queryRunner.manager.save(newImage);
      const uniqueFilename = `/${res.id}.${extension}`;
      res.link = uniqueFilename;

      await queryRunner.manager.save(res);

      if (!fs.existsSync(folder)) {
        fs.mkdirSync(folder);
      }

      const path = `${folder}/${uniqueFilename}`;
      fs.writeFileSync(path, file.buffer);
      console.log(`File uploaded successfully: ${path}`);
      await queryRunner.commitTransaction();
      return res;
    } catch (error) {
      await queryRunner.rollbackTransaction();
      throw error;
    } finally {
      await queryRunner.release();
    }
  }
}
