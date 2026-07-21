import express from "express";
const app = express();
import v1Router from "../routes/v1";

app.use(express.json());
app.use("/api/v1", v1Router);

app.get("/", (req, res) => {
    return res.json({
        "message": "Healthy"
    })
})

app.listen(8000, () => {
    console.log("App is running on port 3000");
})  