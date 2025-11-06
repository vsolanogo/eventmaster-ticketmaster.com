import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, HttpStatus } from '@nestjs/common';
import * as request from 'supertest';
import { AppModule } from '../src/app.module';
import { v4 as uuidv4 } from 'uuid';
import { RegisterDto } from '../src/user/dto/user-register.dto';
import { LoginDto } from '../src/login/dto/login.dto';
import { faker } from '@faker-js/faker';
import * as cookieParser from 'cookie-parser';
import { generateRandomRegisterDto } from '../src/user/user.helpers';

describe('App Tests (e2e)', () => {
  let app: INestApplication;
  let userToken: string;

  beforeEach(async () => {
    const moduleFixture: TestingModule = await Test.createTestingModule({
      imports: [AppModule],
    }).compile();

    app = moduleFixture.createNestApplication();
    app.use(cookieParser());
    await app.init();
  });

  afterEach(async () => {
    await app.close();
  });

  it('/register (POST)', async () => {
    const registerDto: RegisterDto = generateRandomRegisterDto();

    const response = await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.CREATED);

    expect(response.body.id).toBeDefined();
    expect(response.body.email).toBe(registerDto.email);
  });

  it('/register (POST) - User already exists', async () => {
    const registerDto: RegisterDto = generateRandomRegisterDto();

    const response = await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.CREATED);

    const responseForExisting = await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.BAD_REQUEST);
  });

  const successfulLogin = async (): Promise<{
    userDto: RegisterDto;
    userToken: string;
  }> => {
    const registerDto: RegisterDto = generateRandomRegisterDto();

    const responseCreateUser = await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.CREATED);

    const loginDto: LoginDto = {
      email: registerDto.email,
      password: registerDto.password,
    };

    const response = await request(app.getHttpServer())
      .post('/login')
      .send(loginDto)
      .expect(HttpStatus.OK);

    expect(response.header['set-cookie']).toBeDefined();

    userToken = response.header['set-cookie'][0];
    return { userDto: registerDto, userToken };
  };

  it('/login (POST) - Successful login', async () => {
    await successfulLogin();
  });

  it('/login (POST) - Invalid credentials', async () => {
    const loginDto: LoginDto = {
      email: 'nonexistent@example.com',
      password: 'invalidpassword',
    };

    const response = await request(app.getHttpServer())
      .post('/login')
      .send(loginDto)
      .expect(HttpStatus.UNAUTHORIZED);
  });

  it('/user (GET) - Retrieve user information', async () => {
    const { userDto, userToken } = await successfulLogin();

    const res = await request(app.getHttpServer())
      .get('/user')
      .set('cookie', userToken)
      .expect(HttpStatus.OK);

    expect(res.body.id).toBeDefined();
    expect(res.body.email).toBe(userDto.email);

    expect(res.body.password).toBeUndefined();
  });
});
