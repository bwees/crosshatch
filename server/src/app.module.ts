import { DrizzleBetterSQLiteModule } from '@knaadh/nestjs-drizzle-better-sqlite3';
import {
  ArgumentsHost,
  Catch,
  HttpException,
  Logger,
  Module,
} from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import {
  APP_FILTER,
  APP_INTERCEPTOR,
  APP_PIPE,
  BaseExceptionFilter,
} from '@nestjs/core';
import { EventEmitterModule } from '@nestjs/event-emitter';
import {
  ZodSerializationException,
  ZodSerializerInterceptor,
  ZodValidationPipe,
} from 'nestjs-zod';
import { ZodError } from 'zod';
import { controllers } from './controllers';
import { schema } from './db/schema';
import { gateways } from './gateways';
import { proxies } from './proxies';
import { repositories } from './repositories';
import { services } from './services';

const database = DrizzleBetterSQLiteModule.register({
  tag: 'db',
  sqlite3: {
    filename: 'crosshatch.db',
  },
  config: { schema: { ...schema } },
});

@Catch(HttpException)
class HttpExceptionFilter extends BaseExceptionFilter {
  private logger = new Logger(HttpExceptionFilter.name);

  catch(exception: HttpException, host: ArgumentsHost) {
    if (exception instanceof ZodSerializationException) {
      const zodError = exception.getZodError();

      if (zodError instanceof ZodError) {
        this.logger.error(`ZodSerializationException: ${zodError.message}`);
      }
    }

    super.catch(exception, host);
  }
}

@Module({
  imports: [ConfigModule.forRoot(), EventEmitterModule.forRoot(), database],
  controllers: [...controllers],
  providers: [
    ...services,
    ...gateways,
    ...repositories,
    ...proxies,
    {
      provide: APP_PIPE,
      useClass: ZodValidationPipe,
    },
    {
      provide: APP_INTERCEPTOR,
      useClass: ZodSerializerInterceptor,
    },
    {
      provide: APP_FILTER,
      useClass: HttpExceptionFilter,
    },
  ],
})
export class AppModule {}
