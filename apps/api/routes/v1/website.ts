import {Router, type Router as ExpressRouter} from "express";

const websiteRouter: ExpressRouter = Router();

websiteRouter.post("/create-website", (req, res) => {

})

websiteRouter.post("/status/:websiteId", (req, res) => {
    
})

export default websiteRouter;