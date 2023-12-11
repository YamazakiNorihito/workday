import 'reflect-metadata';
import { container } from 'tsyringe';
import { RedisClientType, createClient } from 'redis';
import { FreeeController } from './controllers/freee';
import { SlackHttpApiClient } from './httpClients/slackHttpClient';
import { FreeeAuthenticationService, FreeeService, IFreeeAuthenticationService } from './services/freeeHrService';
import { WeekdayService } from './services/weekdayService';
import { CognitoOAuth2Service, ILoginAuthenticationService } from './services/oauth2Service';
import { LoginController } from './controllers/login';

const redisClient: RedisClientType = createClient({
    url: process.env.REDIS_URL!
});
redisClient.on('error', err => { throw err });
(async () => {
    await redisClient.connect();
})();
// add singleton
container.registerInstance<RedisClientType>("RedisClient", redisClient);

// services
container.register<IFreeeAuthenticationService>("IFreeeAuthenticationService",
    { useValue: new FreeeAuthenticationService(process.env.FREE_CLIENT_ID!, process.env.FREE_CLIENT_SECRET!) });
container.register<SlackHttpApiClient>(SlackHttpApiClient,
    { useValue: new SlackHttpApiClient(process.env.SLACK_TOKEN!) });
container.register<ILoginAuthenticationService>("LoginAuthenticationService",
    { useValue: new CognitoOAuth2Service(process.env.COGNITO_DOMAIN!, process.env.COGNITO_USER_POOL_URL!, process.env.COGNITO_CLIENT_ID!, process.env.COGNITO_CLIENT_SECRET!) });

// controller
const freeeService = container.resolve(FreeeService);
const weekdayService = container.resolve(WeekdayService);
const loginAuthenticationService: ILoginAuthenticationService = container.resolve("LoginAuthenticationService");
container.register<FreeeController>(FreeeController,
    { useValue: new FreeeController(freeeService, weekdayService, process.env.DOMAIN!) });
container.register<LoginController>(LoginController,
    { useValue: new LoginController(loginAuthenticationService, process.env.DOMAIN!) });

export default container;
