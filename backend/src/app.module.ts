import {Module} from '@nestjs/common';
import {AppController} from './app.controller';
import {AppService} from './app.service';
import {SocketModule} from './modules/socket/socket.module';
import { ConverterModule } from './modules/converter/converter.module';

@Module({
  imports: [
    SocketModule,
    ConverterModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {
}
