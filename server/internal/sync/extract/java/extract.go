package java

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/java"
	"go.uber.org/zap"

	"github.com/humanbeeng/lepo/server/internal/sync/extract"
)

type JavaExtractor struct {
	targetTypes           []extract.ChunkType
	targetTypesToRulesDir map[extract.ChunkType]string
	logger                *zap.Logger
}

func NewJavaExtractor(logger *zap.Logger) *JavaExtractor {
	var targetTypes []extract.ChunkType
	targetTypes = append(targetTypes,
		extract.Method,
		extract.Class,
		// TODO: Add interface rule
		extract.Import,
	)

	targetTypesToRulesDir := make(map[extract.ChunkType]string)
	targetTypesToRulesDir[extract.Method] = "internal/sync/extract/java/rules/method.yml"
	targetTypesToRulesDir[extract.Import] = "internal/sync/extract/java/rules/rules/import.yml"
	targetTypesToRulesDir[extract.Constructor] = "internal/sync/extract/java/rules/constructor.yml"
	targetTypesToRulesDir[extract.Import] = "internal/sync/extract/java/rules/import.yml"
	targetTypesToRulesDir[extract.Package] = "internal/sync/extract/java/rules/package.yml"
	targetTypesToRulesDir[extract.Field] = "internal/sync/extract/java/rules/field.yml"

	return &JavaExtractor{
		targetTypes:           targetTypes,
		targetTypesToRulesDir: targetTypesToRulesDir,
		logger:                logger,
	}
}

func (je *JavaExtractor) Extract(file string) ([]extract.Chunk, error) {
	fileinfo, err := os.Stat(file)
	if err != nil {
		je.logger.Error("error:", zap.Error(err))
		return nil, err
	}

	if fileinfo.IsDir() {
		err := fmt.Errorf("error: %v is a directory", file)
		return nil, err
	}

	ext := filepath.Ext(file)
	if ext != ".java" {
		err := fmt.Errorf("error: %v is not a java file", file)
		return nil, err
	}

	lang := java.GetLanguage()
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	// fb, err := os.ReadFile(file) // b has type []byte
	// if err != nil {
	// err := fmt.Errorf("error: Unable to read %v", file)
	// return nil, err
	// }

	fb := `

package com.kyc.report.gst.generator.internal;

import static com.kyc.api.CommonConstants.GST;
import static com.kyc.core.util.KYCConstants.DOT;
import static com.kyc.core.util.KYCConstants.JSON;
import static com.kyc.core.util.KYCConstants.UNDERSCORE;
import static com.kyc.core.util.KYCConstants.XML;


// This is a line comment
class Outer_Demo {
   int num;
   
   @Inject
   public Outer_Demo()  {
   	
   }
   
   Outer_Demo(int num) {
   	this.num = num;
   }
   
   /**
   This is a block comment
   */
   public void hello() {
   
   }

	public void hi(){}

	public void lmao(){}
   
   // inner class
   private class Inner_Demo {
      public void print() throws RuntimeException {
         System.out.println("This is an inner class");
      }
   }
   
  static Outer o = new Outer() {
        void show()
        {
            super.show();
            System.out.println("Demo class");
        }
    };
    
     void outerMethod()
    {
        System.out.println("Outer Method");
        class Inner {
            void innerMethod()
            {
                System.out.println("Inner Method");
            }
        }
 
        Inner y = new Inner();
        y.innerMethod();
    }
   
   // Accessing he inner class from the method within
   void display_Inner() {
      Inner_Demo inner = new Inner_Demo();
      inner.print();
   }
}
	`

	tree, _ := parser.ParseCtx(context.Background(), nil, []byte(fb))
	// queryString := `((class_declaration) @declaration.class)`
	queryString := `
(
    	class_declaration
        body: (
        	class_body(	
            (method_declaration)* @declaration.method
          )  
        )
        
    )* @declaration.class
	`

	q, err := sitter.NewQuery([]byte(queryString), lang)
	if err != nil {
		return nil, fmt.Errorf("Unable to create query")
	}

	rootNode := tree.RootNode()

	qc := sitter.NewQueryCursor()
	qc.Exec(q, rootNode)

	var funcs []*sitter.Node
	f, _ := os.Create("tmp")
	defer f.Close()
	i := 0

	for {
		m, ok := qc.NextMatch()
		i++
		if !ok {
			break
		}

		for _, c := range m.Captures {
			funcs = append(funcs, c.Node)
			f.WriteString(funcName([]byte(fb), c.Node))
		}
	}
	fmt.Println("for ", i)
	return nil, err
}

func funcName(content []byte, n *sitter.Node) string {
	if n == nil {
		return ""
	}
	fmt.Println("Number of named nodes", n.NamedChildCount())
	fmt.Println("Field Name", n.FieldNameForChild(0))

	return fmt.Sprintf("\n\n-Type: %v \nContent\n %v \n----------", n.Type(), n.Content(content))
}
