  import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:async';

void main() {
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'Flutter Demo',
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: Colors.deepPurple),
        useMaterial3: true,
      ),
      home: const MyHomePage(title: 'PixelArt'),
    );
  }
}

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});
  final String title;
  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage> {

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        backgroundColor: Theme.of(context).colorScheme.inversePrimary,
        title: Text(widget.title),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: <Widget>[
            const GenerateImageForm()
          ],
        ),
      ),
    );
  }
}


class GenerateImageForm extends StatefulWidget {
  const GenerateImageForm({super.key});

  @override
  State<StatefulWidget> createState() => _GenerateImageState();
}


class _GenerateImageState extends State<GenerateImageForm> {
  String url = "";
  String id = "";
  final Map<String, TextEditingController> formController = {
    'prompt': TextEditingController()
  };
  
void SubmitPrompt(String prompt) {
  requestId(prompt);
}

Future<String> requestId(String prompt) async {
  print(prompt);
  var url = Uri.https("api.aimless.it","/image");
  // var response = http.post(url, body: {'prompt': prompt});
  var reqBody = {
    "action": 0,
    "params": {
      'prompt': prompt,
      'size': '512x512',
      'n': 1,
      'responseFormat': 'URL',
      'user': 'me'
    }
  };
  var response = await http.post(url, body: json.encode(reqBody), headers: {'Content-Type': 'application/json'});
  print(response.statusCode);
  print(response.body);
  Map<String, dynamic> body = jsonDecode(response.body);
  var id = body['id'];
  print(id);
  if (response.statusCode != 200) {
    throw Exception('Failed to submit prompt');
  }
  setState(() {
    this.id = id;
  });

  Timer.periodic(new Duration(seconds: 5), (Timer timer) async {
    print("Periodic polling for image");
    var statusUrl = Uri.https("api.aimless.it","/status/${id}");
    var statusResponse = await http.get(statusUrl);
    Map<String, dynamic> statusBody = jsonDecode(statusResponse.body);
    if (statusBody['url'] != "") {
      print("Image is ready");
      setState(() {
        this.url = statusBody['url'];
      });
      print("The url is ${this.url}");
      timer.cancel();
    }
  });
  return id;
}

  @override
  Widget build(BuildContext context) {
    return Form(
      child: Column(
        children: [
          TextFormField(
            decoration: const InputDecoration(
              hintText: 'Prompt'
            ),
            validator: (String? value) {
              if (value == null || value.isEmpty){
                return "Enter the prompt for the ai engine";
              }
              return null;
            },
            controller: formController['prompt'],
          ),
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10.0),
            child: ElevatedButton(
              onPressed: () {
                var data = formController['prompt'];
                if (data == null){
                  return;
                }
                SubmitPrompt(data.text);
              },
              child: const Text('Submit'),
            )
          ),
          if (this.url != "") 
            DisplayAIImage(image: this.url)
          
        ],
      )
    );
  }
}


class DisplayAIImage extends StatelessWidget {
  const DisplayAIImage({super.key, required this.image});
  final String image;

  @override
  Widget build(BuildContext context) {
    return Image.network(this.image);
  }
}