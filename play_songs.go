package main


import (
 "context"
 "os"
 "time"


 "github.com/edaniels/golog"
 "github.com/faiface/beep"
 "github.com/faiface/beep/mp3"
 "github.com/faiface/beep/speaker"
 "go.viam.com/rdk/robot/client"
 "go.viam.com/rdk/utils"
 "go.viam.com/utils/rpc"
 "go.viam.com/rdk/services/vision"
)


func initSpeaker(logger golog.Logger) {
   f, err := os.Open("square.mp3")
   if err != nil {
       logger.Fatal(err)
   }
   defer f.Close()


   streamer, format, err := mp3.Decode(f)
   if err != nil {
       logger.Fatal(err)
   }
   defer streamer.Close()


   speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
}


func play(label string, logger golog.Logger) {
   f, err := os.Open(label + ".mp3")
   if err != nil {
       logger.Fatal(err)
   }
   defer f.Close()


   streamer, _, err := mp3.Decode(f)
   if err != nil {
       logger.Fatal(err)
   }
   defer streamer.Close()


   done := make(chan bool)
   speaker.Play(beep.Seq(streamer, beep.Callback(func() {
       done <- true
   })))


   <-done
}


func main() {
 logger := golog.NewDevelopmentLogger("client")
 robot, err := client.New(
    context.Background(),
    "ADDRESS FROM THE VIAM APP",
    logger,
    client.WithDialOptions(rpc.WithEntityCredentials(
    // Replace "<API-KEY-ID>" (including brackets) with your robot's api key id
    "<API-KEY-ID>",
    rpc.Credentials{
        Type:    rpc.CredentialsTypeAPIKey,
        // Replace "<API-KEY>" (including brackets) with your robot's api key
        Payload: "<API-KEY>",
    })),
)
 if err != nil {
    logger.Fatal(err)
 }



 defer robot.Close(context.Background())
  
 visService, err := vision.FromRobot(robot, "shape-classifier")
 if err != nil {
   logger.Error(err)
 }


 initSpeaker(logger)


 for {
   for i := 0; i < 3; i++ {
       visService.ClassificationsFromCamera(context.Background(), "cam",  1, nil)
   }


   classifications, err := visService.ClassificationsFromCamera(context.Background(), "cam", 1, nil)
   if err != nil {
       logger.Fatalf("Could not get classifications: %v", err)
   }
   if len(classifications) > 0 && classifications[0].Score() > 0.7 {
       logger.Info(classifications[0])
       play(classifications[0].Label(), logger)
   }


 }
}

