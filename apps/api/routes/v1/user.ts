import { Router, type Router as ExpressRouter } from "express"

const usersRouter: ExpressRouter = Router();

usersRouter.get("/", (req, res) => {
    res.json({
        message: "Healthy",
    })
})

export default usersRouter;