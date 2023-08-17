  import 'dart:convert';

import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'dart:async';

void main() {
  runApp(const MyApp());
}
class RequestBodyParams {
  RequestBodyParams({this.prompt="", this.size="256x256", this.n=1, this.responseFormat="URL", this.user="1"});
  String prompt;
  String size;
  int n;
  String responseFormat;
  String user;
  Map<String, dynamic> toJson() {
    return {
      'Prompt': prompt,
      'Size': size,
      'N': n,
      'ResponseFormat': responseFormat,
      'User': user
    };
  }
}
class RequestBody {
  RequestBody({this.action=0, required this.params});
  int action;
  RequestBodyParams params;
  Map<String, dynamic> toJson() {
    return {
      'action': action,
      'params': params.toJson()
    };
  }
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

// class PAResponse {
//   urls: List<String>;
// }

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
  RequestBody rb = RequestBody(action: 0, params: RequestBodyParams(prompt: prompt, size: '512x512', n:  1,responseFormat:  'URL', user: 'me'));
  // trying to submit an integer
  // reqBody["params"]["N"] = 1;
  var rbody = json.encode(rb.toJson());
  print(rbody);
  var response = await http.post(url, body: json.encode(rb.toJson()), headers: {'Content-Type': 'application/json'});
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

  Timer.periodic(new Duration(seconds: 3), (Timer timer) async {
    print("Periodic polling for image");
    var statusUrl = Uri.https("api.aimless.it","/status/${id}");
    var statusResponse = await http.get(statusUrl);
    if (statusResponse.statusCode == 204) {
      print("Image is not ready");
      return;
    }
    // PAResponse statusBody = jsonDecode(statusResponse.body);
    if (statusResponse.statusCode == 204) {
      return;
    }
    Map<String, dynamic> statusBody = jsonDecode(statusResponse.body);
    print(statusBody);
    var l = statusBody["urls"] as List<dynamic>;
    if (l != null && l.length > 0) {
      print("Image is ready");
      setState(() {
        this.url = l.first.toString();
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