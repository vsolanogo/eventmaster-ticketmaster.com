import { Test, TestingModule } from '@nestjs/testing';
import { INestApplication, HttpStatus, ValidationPipe } from '@nestjs/common';
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
    app.useGlobalPipes(
      new ValidationPipe({
        transform: true,
        whitelist: true,
      }),
    );
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

  // ==================== AUTHENTICATION EDGE CASES ====================

  it('/user (GET) - Unauthorized without token', async () => {
    await request(app.getHttpServer())
      .get('/user')
      .expect(HttpStatus.FORBIDDEN);
  });

  it('/user (GET) - Unauthorized with invalid token', async () => {
    await request(app.getHttpServer())
      .get('/user')
      .set('cookie', 'SessionID=invalid-token-12345')
      .expect(HttpStatus.NOT_FOUND);
  });

  it('/logout (POST) - Successfully logout', async () => {
    const { userToken } = await successfulLogin();

    const res = await request(app.getHttpServer())
      .post('/logout')
      .set('cookie', userToken)
      .expect(HttpStatus.OK);

    expect(res.header['set-cookie']).toBeDefined();
    expect(res.header['set-cookie'][0]).toContain('SessionID=;');
  });

  // ==================== REGISTRATION VALIDATION ====================

  it('/register (POST) - Invalid email format', async () => {
    const registerDto = {
      email: 'not-an-email',
      password: 'ValidPassword123',
    };

    await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.BAD_REQUEST);
  });

  it('/register (POST) - Missing password', async () => {
    const registerDto = {
      email: 'test@example.com',
    };

    await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.BAD_REQUEST);
  });

  it('/register (POST) - Empty email', async () => {
    const registerDto = {
      email: '',
      password: 'ValidPassword123',
    };

    await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.BAD_REQUEST);
  });

  // ==================== LOGIN EDGE CASES ====================

  it('/login (POST) - Wrong password for existing user', async () => {
    const registerDto = generateRandomRegisterDto();
    
    await request(app.getHttpServer())
      .post('/register')
      .send(registerDto)
      .expect(HttpStatus.CREATED);

    const loginDto: LoginDto = {
      email: registerDto.email,
      password: 'WrongPassword123',
    };

    await request(app.getHttpServer())
      .post('/login')
      .send(loginDto)
      .expect(HttpStatus.UNAUTHORIZED);
  });

  it('/login (POST) - Missing credentials', async () => {
    await request(app.getHttpServer())
      .post('/login')
      .send({})
      .expect(HttpStatus.BAD_REQUEST);
  });

  // ==================== EVENT TESTS ====================

  it('/events (GET) - Get all events without authentication', async () => {
    const res = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    expect(res.body.events).toBeDefined();
    expect(res.body.totalCount).toBeDefined();
    expect(Array.isArray(res.body.events)).toBe(true);
  });

  it('/events (GET) - Pagination works correctly', async () => {
    const res = await request(app.getHttpServer())
      .get('/events?page=1&limit=5')
      .expect(HttpStatus.OK);

    expect(res.body.events).toBeDefined();
    expect(res.body.events.length).toBeLessThanOrEqual(5);
  });

  it('/events (GET) - Sorting by eventDate ASC', async () => {
    const res = await request(app.getHttpServer())
      .get('/events?sortBy=eventDate&sortOrder=ASC')
      .expect(HttpStatus.OK);

    expect(res.body.events).toBeDefined();
    
    if (res.body.events.length > 1) {
      const firstDate = new Date(res.body.events[0].eventDate);
      const secondDate = new Date(res.body.events[1].eventDate);
      expect(firstDate.getTime()).toBeLessThanOrEqual(secondDate.getTime());
    }
  });

  it('/events (GET) - Sorting by eventDate DESC', async () => {
    const res = await request(app.getHttpServer())
      .get('/events?sortBy=eventDate&sortOrder=DESC')
      .expect(HttpStatus.OK);

    expect(res.body.events).toBeDefined();
    
    if (res.body.events.length > 1) {
      const firstDate = new Date(res.body.events[0].eventDate);
      const secondDate = new Date(res.body.events[1].eventDate);
      expect(firstDate.getTime()).toBeGreaterThanOrEqual(secondDate.getTime());
    }
  });

  it('/events/:id (GET) - Get event by ID', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const res = await request(app.getHttpServer())
        .get(`/events/${eventId}`)
        .expect(HttpStatus.OK);

      expect(res.body.id).toBe(eventId);
      expect(res.body.title).toBeDefined();
      expect(res.body.description).toBeDefined();
    }
  });

  it('/events/:id (GET) - Non-existent event returns 404', async () => {
    const fakeId = uuidv4();

    await request(app.getHttpServer())
      .get(`/events/${fakeId}`)
      .expect(HttpStatus.NOT_FOUND);
  });

  it('/events (POST) - Create event requires authentication', async () => {
    const createEventDto = {
      title: 'Test Event',
      description: 'Test Description',
      organizer: 'Test Organizer',
      latitude: 40.7128,
      longitude: -74.0060,
      eventDate: new Date().toISOString(),
      images: [],
    };

    await request(app.getHttpServer())
      .post('/events')
      .send(createEventDto)
      .expect(HttpStatus.FORBIDDEN);
  });

  // ==================== IMAGE UPLOAD TESTS ====================

  it('/image (POST) - Upload requires authentication', async () => {
    await request(app.getHttpServer())
      .post('/image')
      .attach('file', Buffer.from('fake-image-data'), 'test.jpg')
      .expect(HttpStatus.FORBIDDEN);
  });

  // ==================== PARTICIPANT TESTS ====================

  it('/participant/event/:eventId (GET) - Get participants for event', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const res = await request(app.getHttpServer())
        .get(`/participant/event/${eventId}`)
        .expect(HttpStatus.OK);

      expect(Array.isArray(res.body)).toBe(true);
    }
  });

  it('/participant/event/:eventId/registrations-per-day (GET) - Get registration analytics', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const res = await request(app.getHttpServer())
        .get(`/participant/event/${eventId}/registrations-per-day`)
        .expect(HttpStatus.OK);

      expect(Array.isArray(res.body)).toBe(true);
      
      if (res.body.length > 0) {
        expect(res.body[0].date).toBeDefined();
        expect(res.body[0].count).toBeDefined();
      }
    }
  });

  it('/participant/event/:eventId (POST) - Register for event', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const participantDto = {
        fullName: faker.person.fullName(),
        email: faker.internet.email(),
        dateOfBirth: '1990-01-01',
        sourceOfEventDiscovery: 'Social media',
      };

      const res = await request(app.getHttpServer())
        .post(`/participant/event/${eventId}`)
        .send(participantDto)
        .expect(HttpStatus.CREATED);

      expect(res.body.id).toBeDefined();
      expect(res.body.fullName).toBe(participantDto.fullName);
      expect(res.body.email).toBe(participantDto.email);
    }
  });

  it('/participant/event/:eventId (POST) - Multiple registrations allowed', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const participantDto = {
        fullName: faker.person.fullName(),
        email: faker.internet.email(),
        dateOfBirth: '1990-01-01',
        sourceOfEventDiscovery: 'Friends',
      };

      // First registration
      const res1 = await request(app.getHttpServer())
        .post(`/participant/event/${eventId}`)
        .send(participantDto)
        .expect(HttpStatus.CREATED);

      // Second registration with same email is allowed
      const res2 = await request(app.getHttpServer())
        .post(`/participant/event/${eventId}`)
        .send(participantDto)
        .expect(HttpStatus.CREATED);

      expect(res1.body.id).not.toBe(res2.body.id);
      expect(res1.body.email).toBe(res2.body.email);
    }
  });

  it('/participant/event/:eventId (POST) - Invalid participant data fails', async () => {
    const eventsRes = await request(app.getHttpServer())
      .get('/events')
      .expect(HttpStatus.OK);

    if (eventsRes.body.events.length > 0) {
      const eventId = eventsRes.body.events[0].id;

      const invalidParticipantDto = {
        fullName: '',
        email: 'not-an-email',
        dateOfBirth: 'invalid-date',
      };

      await request(app.getHttpServer())
        .post(`/participant/event/${eventId}`)
        .send(invalidParticipantDto)
        .expect(HttpStatus.BAD_REQUEST);
    }
  });

  it('/participant/event/:eventId (GET) - Non-existent event returns 404', async () => {
    const fakeEventId = uuidv4();

    await request(app.getHttpServer())
      .get(`/participant/event/${fakeEventId}`)
      .expect(HttpStatus.NOT_FOUND);
  });
});
