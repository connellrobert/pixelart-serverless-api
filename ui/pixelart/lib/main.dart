  import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

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
      home: const MyHomePage(title: 'Flutter Demo Home Page'),
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
  var url = Uri.https("yesno.wtf","/api");
  // var response = http.post(url, body: {'prompt': prompt});
  var response = await http.get(url);
  print(response.statusCode);
  print(response.body);
  var id = "1-2-3";
  Map<String, dynamic> body = jsonDecode(response.body);
  setState(() {
    this.url = body['image'];
    this.id = id;
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