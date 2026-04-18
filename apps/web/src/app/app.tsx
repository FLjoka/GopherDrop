import { useState, useEffect, useCallback } from 'react';
import axios from 'axios';
import { 
  Upload, 
  Download, 
  Loader2, 
  FileText, 
  Image as ImageIcon,
  HardDrive,
  Trash2,
  Eye
} from 'lucide-react';
import { cn } from '../lib/utils';
import { useToast } from '../hooks/use-toast';
import { Toaster } from '../components/ui/toaster';
import { Button } from '../components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '../components/ui/dialog';

const API_URL = 'http://localhost:8080';

interface FileMetadata {
  ID: number;
  original_name: string;
  size: number;
  content_type: string;
  thumbnail_path: string;
  CreatedAt: string;
}

export function App() {
  const [files, setFiles] = useState<FileMetadata[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [dragActive, setDragActive] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [fileToDelete, setFileToDelete] = useState<number | null>(null);
  const [previewFile, setPreviewFile] = useState<FileMetadata | null>(null);
  const { toast } = useToast();

  const fetchFiles = useCallback(async () => {
    try {
      const response = await axios.get(`${API_URL}/files`);
      setFiles(response.data || []);
    } catch (error) {
      console.error('Error fetching files:', error);
      toast({
        variant: "destructive",
        title: "Error de conexión",
        description: "No se pudieron cargar los archivos del servidor.",
      });
    } finally {
      setIsLoading(false);
    }
  }, [toast]);

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === "dragenter" || e.type === "dragover") {
      setDragActive(true);
    } else if (e.type === "dragleave") {
      setDragActive(false);
    }
  };

  const handleUpload = async (file: File) => {
    setIsUploading(true);
    const formData = new FormData();
    formData.append('file', file);

    try {
      await axios.post(`${API_URL}/upload`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' }
      });
      toast({
        title: "¡Éxito!",
        description: `El archivo ${file.name} se subió correctamente.`,
      });
      fetchFiles();
    } catch (error) {
      console.error('Error uploading file:', error);
      toast({
        variant: "destructive",
        title: "Error de subida",
        description: "Hubo un problema al intentar subir el archivo.",
      });
    } finally {
      setIsUploading(false);
      setDragActive(false);
    }
  };

  const handleDelete = async () => {
    if (fileToDelete === null) return;
    
    try {
      await axios.delete(`${API_URL}/files/${fileToDelete}`);
      setFiles(files.filter(f => f.ID !== fileToDelete));
      toast({
        title: "Archivo eliminado",
        description: "El archivo ha sido borrado físicamente del servidor.",
      });
    } catch (error) {
      console.error('Error deleting file:', error);
      toast({
        variant: "destructive",
        title: "Error al eliminar",
        description: "No se pudo completar la eliminación del archivo.",
      });
    } finally {
      setFileToDelete(null);
    }
  };

  const onDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      handleUpload(e.dataTransfer.files[0]);
    }
  };

  const formatSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  const getFileIcon = (type: string) => {
    if (type.startsWith('image/')) return <ImageIcon className="w-5 h-5 text-blue-500" />;
    return <FileText className="w-5 h-5 text-gray-500" />;
  };

  return (
    <div className="min-h-screen bg-zinc-50 p-8 font-sans text-zinc-900">
      <div className="max-w-5xl mx-auto space-y-8">
        
        {/* Header */}
        <header className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="bg-indigo-600 p-2 rounded-xl">
              <HardDrive className="w-8 h-8 text-white" />
            </div>
            <div>
              <h1 className="text-2xl font-bold tracking-tight">GopherDrop</h1>
              <p className="text-zinc-500 text-sm">Servidor de archivos ultra rápido</p>
            </div>
          </div>
          <div className="flex items-center gap-2 px-3 py-1 bg-green-100 text-green-700 rounded-full text-xs font-medium uppercase tracking-wider">
            <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
            Online
          </div>
        </header>

        {/* Upload Zone */}
        <div 
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={onDrop}
          className={cn(
            "relative group border-2 border-dashed rounded-2xl p-12 transition-all duration-200 flex flex-col items-center justify-center gap-4 cursor-pointer",
            dragActive ? "border-indigo-500 bg-indigo-50" : "border-zinc-200 bg-white hover:border-zinc-300"
          )}
        >
          <input 
            type="file" 
            className="absolute inset-0 opacity-0 cursor-pointer"
            onChange={(e) => e.target.files?.[0] && handleUpload(e.target.files[0])}
          />
          <div className={cn(
            "p-4 rounded-full transition-colors",
            dragActive ? "bg-indigo-100 text-indigo-600" : "bg-zinc-100 text-zinc-500 group-hover:bg-zinc-200"
          )}>
            {isUploading ? <Loader2 className="w-8 h-8 animate-spin" /> : <Upload className="w-8 h-8" />}
          </div>
          <div className="text-center">
            <p className="font-semibold">{isUploading ? 'Subiendo archivo...' : 'Suelte archivos aquí para subir'}</p>
            <p className="text-sm text-zinc-500 mt-1">O haz clic para seleccionar manualmente</p>
          </div>
        </div>

        {/* File Table */}
        <div className="bg-white rounded-2xl border border-zinc-200 overflow-hidden shadow-sm">
          <div className="p-6 border-b border-zinc-100 flex items-center justify-between">
            <h2 className="font-semibold text-lg">Archivos Recientes</h2>
            <span className="text-xs font-medium text-zinc-400 uppercase tracking-widest">{files.length} archivos</span>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-zinc-50/50">
                  <th className="px-6 py-4 text-sm font-semibold text-zinc-600">Nombre</th>
                  <th className="px-6 py-4 text-sm font-semibold text-zinc-600">Tamaño</th>
                  <th className="px-6 py-4 text-sm font-semibold text-zinc-600">Fecha</th>
                  <th className="px-6 py-4 text-sm font-semibold text-zinc-600 text-right">Acción</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-zinc-100">
                {isLoading ? (
                   <tr>
                    <td colSpan={4} className="px-6 py-12 text-center text-zinc-400">
                      <Loader2 className="w-6 h-6 animate-spin mx-auto mb-2" />
                      Cargando archivos...
                    </td>
                  </tr>
                ) : files.length === 0 ? (
                  <tr>
                    <td colSpan={4} className="px-6 py-12 text-center text-zinc-400">
                      No hay archivos subidos todavía.
                    </td>
                  </tr>
                ) : (
                  files.map((file) => (
                    <tr key={file.ID} className="group hover:bg-zinc-50 transition-colors">
                      <td className="px-6 py-4">
                        <div className="flex items-center gap-3">
                          {file.thumbnail_path ? (
                            <img 
                              src={`${API_URL}/thumbnails/${file.ID}`} 
                              alt={file.original_name}
                              className="w-10 h-10 rounded-md object-cover border border-zinc-200 shadow-sm"
                            />
                          ) : (
                            <div className="w-10 h-10 rounded-md bg-zinc-100 flex items-center justify-center border border-zinc-200">
                              {getFileIcon(file.content_type)}
                            </div>
                          )}
                          <span className="font-medium truncate max-w-[240px]" title={file.original_name}>
                            {file.original_name}
                          </span>
                        </div>
                      </td>
                      <td className="px-6 py-4 text-sm text-zinc-500">
                        {formatSize(file.size)}
                      </td>
                      <td className="px-6 py-4 text-sm text-zinc-500">
                        {new Date(file.CreatedAt).toLocaleDateString()}
                      </td>
                      <td className="px-6 py-4 text-right">
                        <div className="flex items-center justify-end gap-2">
                          {file.content_type.startsWith('image/') && (
                            <button
                              onClick={() => setPreviewFile(file)}
                              className="inline-flex items-center justify-center p-2 rounded-lg bg-zinc-100 text-zinc-600 hover:bg-indigo-600 hover:text-white transition-all shadow-sm group/btn"
                              title="Previsualizar"
                            >
                              <Eye className="w-4 h-4" />
                            </button>
                          )}
                          <a 
                            href={`${API_URL}/download/${file.ID}`}
                            download={file.original_name}
                            className="inline-flex items-center justify-center p-2 rounded-lg bg-zinc-100 text-zinc-600 hover:bg-indigo-600 hover:text-white transition-all shadow-sm group/btn"
                            title="Descargar"
                          >
                            <Download className="w-4 h-4" />
                          </a>
                          <button
                            onClick={() => setFileToDelete(file.ID)}
                            className="inline-flex items-center justify-center p-2 rounded-lg bg-zinc-100 text-zinc-600 hover:bg-red-600 hover:text-white transition-all shadow-sm group/btn"
                            title="Eliminar"
                          >
                            <Trash2 className="w-4 h-4" />
                          </button>
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Confirmation Dialog */}
      <Dialog open={fileToDelete !== null} onOpenChange={(open) => !open && setFileToDelete(null)}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>¿Confirmar eliminación?</DialogTitle>
            <DialogDescription>
              Esta acción borrará permanentemente el archivo del servidor y su miniatura asociada. Esta acción no se puede deshacer.
            </DialogDescription>
          </DialogHeader>
          <DialogFooter className="sm:justify-end gap-2">
            <Button variant="outline" onClick={() => setFileToDelete(null)}>
              Cancelar
            </Button>
            <Button variant="destructive" onClick={handleDelete}>
              Eliminar Archivo
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* Preview Dialog */}
      <Dialog open={previewFile !== null} onOpenChange={(open) => !open && setPreviewFile(null)}>
        <DialogContent className="max-w-4xl p-0 overflow-hidden bg-transparent border-none shadow-none">
          <DialogTitle className="sr-only">Vista previa de {previewFile?.original_name}</DialogTitle>
          {previewFile && (
            <div className="relative flex items-center justify-center p-4">
              <img 
                src={`${API_URL}/download/${previewFile.ID}`} 
                alt={previewFile.original_name}
                className="max-w-full max-h-[85vh] rounded-lg shadow-2xl object-contain bg-zinc-900"
              />
              <div className="absolute bottom-10 left-1/2 -translate-x-1/2 bg-black/60 backdrop-blur-md text-white px-4 py-2 rounded-full text-sm font-medium border border-white/10 shadow-lg">
                {previewFile.original_name} • {formatSize(previewFile.size)}
              </div>
            </div>
          )}
        </DialogContent>
      </Dialog>

      <Toaster />
    </div>
  );
}

export default App;
