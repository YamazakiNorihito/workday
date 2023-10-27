import express from 'express';
import * as topController from '../controllers/top';

const router = express.Router();

router.get('/', topController.get);

export default router;
