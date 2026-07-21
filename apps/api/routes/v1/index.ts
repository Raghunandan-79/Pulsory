import Router from "express";
import usersRouter from "./user";
import websiteRouter from "./website";

const router = Router();

router.use("/user", usersRouter);
router.use("/website", websiteRouter);

export default router;