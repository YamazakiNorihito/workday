import express from 'express';
import topRoutes from './routes/topRoutes';
import { isProduction } from './environment';
import dotenv from 'dotenv';

dotenv.config();

const app = express();

app.use(express.json());
app.use('/', topRoutes);

app.listen(3000, () => {
    console.log('Server is running on port 3000');
});

export default app;