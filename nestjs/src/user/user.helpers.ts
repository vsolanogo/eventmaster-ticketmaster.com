import { randEmail, randPassword, randFullName } from '@ngneat/falso';
import { RegisterDto } from './dto/user-register.dto';

export const generateRandomRegisterDto = (): RegisterDto => ({
  email: randEmail(),
  password: randPassword(),
});
